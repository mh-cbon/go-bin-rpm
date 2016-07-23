package rpm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/mattn/go-zglob"
	"github.com/mh-cbon/go-bin-rpm/stringexec"
	"github.com/mh-cbon/verbose"
)

var logger = verbose.Auto()

type Package struct {
	Name          string            `json:"name"`
	Version       string            `json:"version,omitempty"`
	Arch          string            `json:"arch,omitempty"`
	Release       string            `json:"release,omitempty"`
	Group         string            `json:"group,omitempty"`
	License       string            `json:"license,omitempty"`
	Url           string            `json:"url,omitempty"`
	Summary       string            `json:"summary,omitempty"`
	Description   string            `json:"description,omitempty"`
	ChangelogFile string            `json:"changelog-file,omitempty"`
	ChangelogCmd  string            `json:"changelog-cmd,omitempty"`
	Files         []FileInstruction `json:"files,omitempty"`
	PreInst       string            `json:"preinst,omitempty"`
	PostInst      string            `json:"postinst,omitempty"`
	PreRm         string            `json:"prerm,omitempty"`
	PostRm        string            `json:"postrm,omitempty"`
	Verify        string            `json:"verify,omitempty"`
	BuildRequires []string          `json:"build-requires,omitempty"`
	Requires      []string          `json:"requires,omitempty"`
	Provides      []string          `json:"provides,omitempty"`
	Conflicts     []string          `json:"conflicts,omitempty"`
	Envs          map[string]string `json:"envs,omitempty"`
	Menus         []Menu            `json:"menus"`
}

type FileInstruction struct {
	From string `json:"from, omitempty"`
	To   string `json:"to, omitempty"`
	Base string `json:"base, omitempty"`
	Type string `json:"type, omitempty"`
}

type Menu struct {
	Name            string `json:"name"`           // Name of the shortcut
	GenericName     string `json:"generic-name"`   //
	Exec            string `json:"exec"`           // Exec command
	Icon            string `json:"icon"`           // Path to the installed icon
	Type            string `json:"type"`           // Type of shortcut
	StartupNotify   bool   `json:"startup-notify"` // yes/no
	Terminal        bool   `json:"terminal"`       // yes/no
	DBusActivatable bool   `json:"dbus-activable"` // yes/no
	NoDisplay       bool   `json:"no-display"`     // yes/no
	Keywords        string `json:"keywords"`       // ; separated list
	OnlyShowIn      string `json:"only-show-in"`   // ; separated list
	Categories      string `json:"categories"`     // ; separated list
	MimeType        string `json:"mime-type"`      // ; separated list
}

func (p *Package) Load(file string) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		m := fmt.Sprintf("json file '%s' does not exist: %s", file, err.Error())
		return errors.New(m)
	}
	byt, err := ioutil.ReadFile(file)
	if err != nil {
		m := fmt.Sprintf("error occured while reading file '%s': %s", file, err.Error())
		return errors.New(m)
	}
	if err := json.Unmarshal(byt, p); err != nil {
		m := fmt.Sprintf("Invalid json file '%s': %s", file, err.Error())
		return errors.New(m)
	}
	return nil
}

func (p *Package) Normalize(arch string, version string) error {

	tokens := make(map[string]string)
	tokens["!version!"] = version
	tokens["!arch!"] = arch
	tokens["!name!"] = p.Name

	p.Version = replaceTokens(p.Version, tokens)
	p.Arch = replaceTokens(p.Arch, tokens)
	p.Url = replaceTokens(p.Url, tokens)
	p.Summary = replaceTokens(p.Summary, tokens)
	p.Description = replaceTokens(p.Description, tokens)
	p.ChangelogFile = replaceTokens(p.ChangelogFile, tokens)
	p.ChangelogCmd = replaceTokens(p.ChangelogCmd, tokens)

	if p.Release == "" {
		p.Release = "1"
	}
	if p.Version == "" {
		p.Version = version
	}
	if p.Arch == "" {
		p.Arch = arch
	}
	for i, v := range p.Files {
		p.Files[i].From = replaceTokens(v.From, tokens)
		p.Files[i].Base = replaceTokens(v.Base, tokens)
		p.Files[i].To = replaceTokens(v.To, tokens)
	}
	logger.Printf("Arch=%s\n", p.Arch)
	logger.Printf("Version=%s\n", p.Version)
	logger.Printf("Release=%s\n", p.Release)
	logger.Printf("Url=%s\n", p.Url)
	logger.Printf("Summary=%s\n", p.Summary)
	logger.Printf("Description=%s\n", p.Description)
	logger.Printf("ChangelogFile=%s\n", p.ChangelogFile)
	logger.Printf("ChangelogCmd=%s\n", p.ChangelogCmd)

	shortcuts, err := p.WriteShortcutFiles()
	if err != nil {
		return err
	}
	logger.Printf("shortcuts=%s\n", shortcuts)
	for _, shortcut := range shortcuts {
		sc := FileInstruction{}
		sc.From = shortcut
		sc.To = fmt.Sprintf("%%{_datadir}/applications/")
		sc.Base = filepath.Dir(shortcut)
		p.Files = append(p.Files, sc)
		logger.Printf("Added menu shortcut File=%q\n", sc)
	}
	for _, menu := range p.Menus {
		sc := FileInstruction{}
		f, err := filepath.Abs(menu.Icon)
		if err != nil {
			return err
		}
		sc.From = f
		sc.To = fmt.Sprintf("%%{_datadir}/pixmaps/")
		sc.Base = filepath.Dir(f)
		p.Files = append(p.Files, sc)
		logger.Printf("Added menu icon File=%q\n", sc)

		// desktop-file-utils is super picky.
		menu.Categories = strings.TrimSuffix(menu.Categories, ";")
		menu.Keywords = strings.TrimSuffix(menu.Keywords, ";")
		if menu.Categories != "" {
			menu.Categories += ";"
		}
		if menu.Keywords != "" {
			menu.Keywords += ";"
		}
	}

	if len(p.Menus) > 0 {
		if contains(p.BuildRequires, "desktop-file-utils") == false {
			p.BuildRequires = append(p.BuildRequires, "desktop-file-utils")
		}
	}

	if len(p.Envs) > 0 {
		envFile, err := p.WriteEnvFile()
		if err != nil {
			return err
		}
		sc := FileInstruction{}
		sc.From = envFile
		sc.To = fmt.Sprintf("%%{_sysconfdir}/profile.d/")
		sc.Base = filepath.Dir(envFile)
		logger.Printf("Added env File=%q\n", sc)
		p.Files = append(p.Files, sc)
	}
	logger.Printf("p.Envs=%s\n", p.Envs)
	logger.Printf("p.Requires=%s\n", p.Requires)
	logger.Printf("p.BuildRequires=%s\n", p.BuildRequires)
	return nil
}

func replaceTokens(in string, tokens map[string]string) string {
	for token, v := range tokens {
		in = strings.Replace(in, token, v, -1)
	}
	return in
}

func (p *Package) InitializeBuildArea(buildAreaPath string) error {
	paths := make([]string, 0)
	paths = append(paths, filepath.Join(buildAreaPath, "BUILD"))
	paths = append(paths, filepath.Join(buildAreaPath, "RPMS"))
	paths = append(paths, filepath.Join(buildAreaPath, "SOURCES"))
	paths = append(paths, filepath.Join(buildAreaPath, "SPECS"))
	paths = append(paths, filepath.Join(buildAreaPath, "SRPMS"))
	paths = append(paths, filepath.Join(buildAreaPath, "RPMS", "i386"))
	paths = append(paths, filepath.Join(buildAreaPath, "RPMS", "amd64"))

	for _, p := range paths {
		if err := os.MkdirAll(p, 0755); err != nil {
			return err
		}
	}
	return nil
}

func (p *Package) WriteSpecFile(sourceDir string, buildAreaPath string) error {
	spec, err := p.GenerateSpecFile(sourceDir)
	if err != nil {
		return err
	}
	path := filepath.Join(buildAreaPath, "SPECS", p.Name+".spec")
	return ioutil.WriteFile(path, []byte(spec), 0644)
}

func (p *Package) RunBuild(buildAreaPath string, output string) error {
	path := filepath.Join(buildAreaPath, "SPECS", p.Name+".spec")
	def := "_topdir " + buildAreaPath
	arch := p.Arch
	if arch == "386" {
		arch = "i386"
	}
	if arch == "amd64" {
		arch = "x86_64"
	}
	args := []string{"--target", arch, "-bb", path, "--define", def}
	logger.Printf("%s %s\n", "rpmbuild", args)
	oCmd := exec.Command("rpmbuild", args...)
	oCmd.Stdout = os.Stdout
	oCmd.Stderr = os.Stderr
	if err := oCmd.Run(); err != nil {
		return err
	}
	pkg := fmt.Sprintf("%s/RPMS/%s/%s-%s-%s.%s.rpm", buildAreaPath, arch, p.Name, p.Version, p.Release, arch)
	return cp(output, pkg)
}

func (p *Package) GenerateSpecFile(sourceDir string) (string, error) {
	spec := ""

	// Version field of the spec file must not
	// contain non numeric characters,
	// see https://fedoraproject.org/wiki/Packaging:Naming?rd=Packaging:NamingGuidelines#Version_Tag
	// the prerelease stuff is moved into Release field
	v, err := semver.NewVersion(p.Version)
	if err != nil {
		return "", err
	}
	okVersion := strings.Replace(v.String(), v.Prerelease(), "", -1)
	preRelease := p.Release
	if v.Prerelease() != "" {
		preRelease = v.Prerelease() + "." + preRelease
	}

	if p.Name != "" {
		spec += fmt.Sprintf("Name: %s\n", p.Name)
	}
	if p.Version != "" {
		spec += fmt.Sprintf("Version: %s\n", okVersion)
	}
	if p.Release != "" {
		spec += fmt.Sprintf("Release: %s\n", preRelease)
	}
	if p.Group != "" {
		spec += fmt.Sprintf("Group: %s\n", p.Group)
	}
	if p.License != "" {
		spec += fmt.Sprintf("License: %s\n", p.License)
	}
	if p.Url != "" {
		spec += fmt.Sprintf("Url: %s\n", p.Url)
	}
	if p.Summary != "" {
		spec += fmt.Sprintf("Summary: %s\n", p.Summary)
	}
	if len(p.BuildRequires) > 0 {
		spec += fmt.Sprintf("\nBuildRequires:%s\n", strings.Join(p.BuildRequires, ", "))
	}
	if len(p.Requires) > 0 {
		spec += fmt.Sprintf("\nRequires:%s\n", strings.Join(p.Requires, ", "))
	}
	if len(p.Provides) > 0 {
		spec += fmt.Sprintf("\nProvides:%s\n", strings.Join(p.Provides, ", "))
	}
	if len(p.Conflicts) > 0 {
		spec += fmt.Sprintf("\nConflicts:%s\n", strings.Join(p.Conflicts, ", "))
	}
	if p.Description != "" {
		spec += fmt.Sprintf("\n%%description\n%s\n", p.Description)
	}
	spec += fmt.Sprintf("\n%%prep\n")
	spec += fmt.Sprintf("\n%%build\n")
	spec += fmt.Sprintf("\n%%install\n")
	if install, err := p.GenerateInstallSection(sourceDir); err != nil {
		return "", err
	} else {
		spec += fmt.Sprintf("%s\n", install)
	}
	spec += fmt.Sprintf("\n%%files\n")
	if files, err := p.GenerateFilesSection(sourceDir); err != nil {
		return "", err
	} else {
		spec += fmt.Sprintf("%s\n", files)
	}
	spec += fmt.Sprintf("\n%%clean\n")
	shortcutInstall := "\n"
	for _, menu := range p.Menus {
		shortcutInstall += fmt.Sprintf("desktop-file-install --vendor='' ")
		shortcutInstall += fmt.Sprintf("--dir=%%{buildroot}%%{_datadir}/applications/%s ", p.Name)
		shortcutInstall += fmt.Sprintf("%%{buildroot}/%%{_datadir}/applications/%s.desktop", menu.Name)
		shortcutInstall += "\n"
	}
	shortcutInstall = strings.TrimSpace(shortcutInstall)
	if content := readFile(p.PreInst); content != "" {
		spec += fmt.Sprintf("\n%%pre\n%s\n", content)
	}
	if content := readFile(p.PostInst); content != "" {
		spec += fmt.Sprintf("\n%%post\n%s\n", content+shortcutInstall)
	} else if shortcutInstall != "" {
		spec += fmt.Sprintf("\n%%post\n%s\n", shortcutInstall)
	}
	if content := readFile(p.PreRm); content != "" {
		spec += fmt.Sprintf("\n%%preun\n%s\n", content)
	}
	if content := readFile(p.PostRm); content != "" {
		spec += fmt.Sprintf("\n%%postun\n%s\n", content)
	}
	if content := readFile(p.Verify); content != "" {
		spec += fmt.Sprintf("\n%%verifyscript\n%s\n", content)
	}
	spec += fmt.Sprintf("\n%%changelog\n")
	if content, err := p.GetChangelogContent(); err != nil {
		return "", err
	} else {
		spec += fmt.Sprintf("%s\n", content)
	}

	return spec, nil
}

func (p *Package) GenerateInstallSection(sourceDir string) (string, error) {
	var err error
	allDirs := make([]string, 0)
	allFiles := make([]string, 0)
	if sourceDir, err = filepath.Abs(sourceDir); err != nil {
		return "", err
	}
	for i, fileInst := range p.Files {

		if fileInst.From == "" {
			logger.Printf("Skipped p.Files[%d] %q", i, fileInst)
			continue
		}

		from := fileInst.From
		to := fileInst.To
		base := fileInst.Base

		if filepath.IsAbs(from) == false {
			from = filepath.Join(sourceDir, from)
		}
		if filepath.IsAbs(base) == false {
			base = filepath.Join(sourceDir, base)
		}

		logger.Printf("fileInst.From=%q\n", from)
		logger.Printf("fileInst.To=%q\n", to)
		logger.Printf("fileInst.Base=%q\n", base)

		items, err := zglob.Glob(from)
		if err != nil {
			m := fmt.Sprintf("Could not glob files source '%s': %s", from, err.Error())
			return "", errors.New(m)
		}
		logger.Printf("items=%q\n", items)
		for _, item := range items {
			n := item
			if len(item) >= len(base) && item[0:len(base)] == base {
				n = item[len(base):]
			}
			n = filepath.Join("%{buildroot}", to, n)
			dir := fmt.Sprintf("mkdir -p %s\n", filepath.Dir(n))
			if contains(allDirs, dir) == false {
				allDirs = append(allDirs, dir)
			}
			if s, err := os.Stat(item); err != nil {
				return "", err
			} else if s.IsDir() == false {
				file := fmt.Sprintf("cp %s %s\n", item, filepath.Dir(n))
				if contains(allFiles, file) == false {
					allFiles = append(allFiles, file)
				}
			}
		}
	}

	content := ""
	for _, d := range allDirs {
		content += d
	}
	for _, d := range allFiles {
		content += d
	}

	logger.Printf("content=\n%s\n", content)

	return content, nil
}

func (p *Package) GenerateFilesSection(sourceDir string) (string, error) {
	var err error
	content := ""
	allItems := make([]fileItem, 0)

	if sourceDir, err = filepath.Abs(sourceDir); err != nil {
		return "", err
	}

	for _, fileInst := range p.Files {
		from := fileInst.From
		to := fileInst.To
		base := fileInst.Base
		ftype := fileInst.Type

		if ftype != "" {
			ftype = " " + ftype
		}

		if from == "" {
			content += fmt.Sprintf("%s\n", ftype)
			continue
		}

		if filepath.IsAbs(from) == false {
			from = filepath.Join(sourceDir, from)
		}
		if filepath.IsAbs(base) == false {
			base = filepath.Join(sourceDir, base)
		}

		logger.Printf("fileInst.From=%q\n", from)
		logger.Printf("fileInst.To=%q\n", to)
		logger.Printf("fileInst.Base=%q\n", base)
		logger.Printf("fileInst.Type=%q\n", ftype)

		items, err := zglob.Glob(from)
		if err != nil {
			m := fmt.Sprintf("Could not glob files source '%s': %s", from, err.Error())
			return "", errors.New(m)
		}
		logger.Printf("items=%q\n", items)
		for _, item := range items {
			n := item
			if len(item) >= len(base) && item[0:len(base)] == base {
				n = item[len(base):]
			}
			n = filepath.Join(to, n)
			if fileItems(allItems).contains(n) == false {
				allItems = append(allItems, fileItem{n, ftype})
			}
		}
	}

	for _, item := range allItems {
		content += fmt.Sprintf("%s%s\n", item.Type, item.Path)
	}

	logger.Printf("content=\n%s\n", content)

	return content, nil
}

func (p *Package) GetChangelogContent() (string, error) {
	var err error
	var c []byte
	var wd string
	var cmd *exec.Cmd
	if p.ChangelogFile != "" {
		if c, err = ioutil.ReadFile(p.ChangelogFile); err == nil {
			return string(c), nil
		}
	} else if p.ChangelogCmd != "" {
		wd, err = os.Getwd()
		if err == nil {
			cmd, err = stringexec.Command(wd, p.ChangelogCmd)
			if err == nil {
				cmd.Stdout = nil
				c, err = cmd.Output()
				if err == nil {
					return string(c), nil
				}
			}
		}
	}
	return "", err
}

func (p *Package) WriteShortcutFiles() ([]string, error) {

	files := make([]string, 0)

	tpmDir, err := ioutil.TempDir("", "rpm-desktops")
	if err != nil {
		return files, err
	}

	for _, m := range p.Menus {
		s := ""

		if m.Name != "" {
			s += fmt.Sprintf("%s=%s\n", "Name", m.Name)
		}

		if m.GenericName != "" {
			s += fmt.Sprintf("%s=%s\n", "GenericName", m.GenericName)
		}

		if m.Exec != "" {
			s += fmt.Sprintf("%s=%s\n", "Exec", m.Exec)
		}

		if m.Icon != "" {
			s += fmt.Sprintf("%s=%s\n", "Icon", "/usr/share/pixmaps/"+filepath.Base(m.Icon))
		}

		if m.Type != "" {
			s += fmt.Sprintf("%s=%s\n", "Type", m.Type)
		}

		if m.Categories != "" {
			s += fmt.Sprintf("%s=%s\n", "Categories", m.Categories+";")
		}

		if m.MimeType != "" {
			s += fmt.Sprintf("%s=%s\n", "MimeType", m.MimeType)
		}

		if m.OnlyShowIn != "" {
			s += fmt.Sprintf("%s=%s\n", "OnlyShowIn", m.OnlyShowIn)
		}

		if m.Keywords != "" {
			s += fmt.Sprintf("%s=%s\n", "Keywords", m.Keywords+";")
		}

		if s != "" {

			if m.StartupNotify {
				s += fmt.Sprintf("%s=%s\n", "StartupNotify", "true")
			} else {
				s += fmt.Sprintf("%s=%s\n", "StartupNotify", "false")
			}

			if m.DBusActivatable {
				s += fmt.Sprintf("%s=%s\n", "DBusActivatable", "true")
			} else {
				s += fmt.Sprintf("%s=%s\n", "DBusActivatable", "false")
			}

			if m.NoDisplay {
				s += fmt.Sprintf("%s=%s\n", "NoDisplay", "true")
			} else {
				s += fmt.Sprintf("%s=%s\n", "NoDisplay", "false")
			}

			if m.Terminal {
				s += fmt.Sprintf("%s=%s\n", "Terminal", "true")
			} else {
				s += fmt.Sprintf("%s=%s\n", "Terminal", "false")
			}

			s = "[Desktop Entry]\n" + s

			file := filepath.Join(tpmDir, m.Name+".desktop")

			files = append(files, file)

			if err := ioutil.WriteFile(file, []byte(s), 0644); err != nil {
				return files, err
			}
		}
	}

	return files, nil
}

func (p *Package) WriteEnvFile() (string, error) {

	file := ""

	tpmDir, err := ioutil.TempDir("", "rpm-envs")
	if err != nil {
		return file, err
	}

	file = filepath.Join(tpmDir, p.Name+".sh")

	content := "#!/bin/bash\n\n"
	for k, v := range p.Envs {
		content += fmt.Sprintf("%s=%s\n", k, v)
	}
	content += fmt.Sprint("\n")
	for k, _ := range p.Envs {
		content += fmt.Sprintf("export %s\n", k)
	}
	return file, ioutil.WriteFile(file, []byte(content), 0644)
}

type fileItem struct {
	Path string
	Type string
}

type fileItems []fileItem

func (f fileItems) contains(path string) bool {
	for _, item := range f {
		if item.Path == path {
			return true
		}
	}
	return false
}

func contains(l []string, v string) bool {
	for _, vv := range l {
		if vv == v {
			return true
		}
	}
	return false
}

func cp(dst, src string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}

func readFile(src string) string {
	if c, err := ioutil.ReadFile(src); err != nil {
		return ""
	} else {
		return string(c)
	}
}
