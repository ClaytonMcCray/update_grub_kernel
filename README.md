# Overview
This is a tool for selecting a default kernel to boot into.

## Usage:

`sudo update_grub_kernel`

Then select the kernel you wish to boot into. It works by
parsing grub.cfg for possible options, presenting them to the
user, then writing to `GRUB_DEFAULT` in /etc/default/grub.

By default it will run `update-grub2` on it's own.

Symlink the installed binary to somewhere on root's $PATH,
like /sbin.

