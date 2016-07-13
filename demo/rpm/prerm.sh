echo "prerm"

# non standard stuff, stop the service asap
systemctl stop hello.service

# unregister the service
%systemd_preun hello.service
