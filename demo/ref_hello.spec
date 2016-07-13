Name: software
Version: 0.23
Release: 1
Summary: Custom software to run enterprise servers.

Group: Applications/Internet
License: GPLv2
URL: http://meinit.nl/

%description
This software runs all enterprise software as a daemon. It's been developed by Me in IT consultancy.

%prep

%build

%install
mkdir -p %{buildroot}/%{_bindir}/
mkdir -p %{buildroot}/%{_docdir}/%{name}-%{version}-%{release}/
cp /vagrant/hello %{buildroot}/%{_bindir}/hello
cp /vagrant/README.md %{buildroot}/%{_docdir}/%{name}-%{version}-%{release}/README.md

%files
%defattr(-,root,root)
%doc %{_docdir}/%{name}-%{version}-%{release}/README.md
%{_bindir}/hello

%clean

%changelog
* Tue Jun 15 2010 Robert de Bock <robert@meinit.nl> - 0.23-1
- Initial build
