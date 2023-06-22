# Systemd

The `.service` file in this directory should be copied to the `/etc/systemd/system/` folder.
- It expects the binary file to be installed at `/usr/sbin/fail2ban_exporter`.
- It expects a user named `fail2ban_exporter` to exist. This user should not have a shell or any special privileges aside from read-access to the fail2ban socket file.

The `ExecStart` line can be modified to add any custom CLI flags.
