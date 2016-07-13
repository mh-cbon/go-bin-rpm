echo "postinst"

# add a group
/usr/bin/getent group hello || /usr/sbin/groupadd -r hello
# add an user
/usr/bin/getent passwd hello || /usr/sbin/useradd hello -g hello --system --no-create-home --home /nonexistent

# register service
%systemd_post hello.service

# non standard stuff, start the service asap
systemctl start hello.service
