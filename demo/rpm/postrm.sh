echo "postrm"

# not sure what this does...
%systemd_postun_with_restart hello.service
