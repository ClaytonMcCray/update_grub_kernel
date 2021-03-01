package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestCutFront(t *testing.T) {
	initial := "152:menuentry 'Ubuntu, with Linux 5.11.0' --class ubuntu --class gnu-linux --class gnu --class os $menuentry_id_option 'gnulinux-5.11.0-advanced-b70cb823-9505-4ab6-bc0a-ca359515bf03' {"

	expected := "'gnulinux-5.11.0-advanced-b70cb823-9505-4ab6-bc0a-ca359515bf03' {"
	actual, err := cutFront(initial)

	if err != nil {
		t.Errorf("expected error to be nil, got %s", err)
	}

	if actual != expected {
		t.Errorf("expected = %s, actual = %s", expected, actual)
	}

	initial = ""
	_, err = cutFront(initial)

	if err == nil {
		t.Errorf("expected error to be nil, got %s", err)
	}
}

func TestCutRear(t *testing.T) {
	initial := "152:menuentry 'Ubuntu, with Linux 5.11.0' --class ubuntu --class gnu-linux --class gnu --class os $menuentry_id_option 'gnulinux-5.11.0-advanced-b70cb823-9505-4ab6-bc0a-ca359515bf03' {"

	expected := "152:menuentry 'Ubuntu, with Linux 5.11.0' --class ubuntu --class gnu-linux --class gnu --class os $menuentry_id_option 'gnulinux-5.11.0-advanced-b70cb823-9505-4ab6-bc0a-ca359515bf03'"

	actual, err := cutRear(initial)

	if err != nil {
		t.Errorf("got error %s, expected nil", err)
	}

	if actual != expected {
		t.Errorf("expected %s, actual = %s", expected, actual)
	}

	initial = ""
	_, err = cutRear(initial)

	if err == nil {
		t.Errorf("expected error to be nil, got %s", err)
	}
}

func TestProcess(t *testing.T) {
	initial := "152:menuentry 'Ubuntu, with Linux 5.11.0' --class ubuntu --class gnu-linux --class gnu --class os $menuentry_id_option 'gnulinux-5.11.0-advanced-b70cb823-9505-4ab6-bc0a-ca359515bf03' {"

	expected := "'gnulinux-5.11.0-advanced-b70cb823-9505-4ab6-bc0a-ca359515bf03'"
	actual, err := process(initial)

	if err != nil {
		t.Errorf("expected nil error, got %s", err)
	}

	if actual != expected {
		t.Errorf("expected = %s, actual = %s", expected, actual)
	}
}

func TestFindKernels(t *testing.T) {
	expected := []string{"\"gnulinux-5.11.0-advanced-b70cb823-9505-4ab6-bc0a-ca359515bf03\"",
		"\"gnulinux-5.11.0-recovery-b70cb823-9505-4ab6-bc0a-ca359515bf03\"",
		"\"gnulinux-5.11.0.old-advanced-b70cb823-9505-4ab6-bc0a-ca359515bf03\""}

	actual := findKernels(examplegrubcfg)

	for i, k := range actual {
		if k != expected[i] {
			t.Errorf("got %s, expected %s", k, expected[i])
		}
	}
}

func TestReadReducedDefaults(t *testing.T) {
	expected := strings.Split(examplegrubdefaultsnolines, "\n")
	actual, err := readReducedDefaults(bytes.NewReader([]byte(examplegrubdefaults)))

	if err != nil {
		t.Errorf("got err = %s, expected nil", err)
	}

	for i, a := range actual {
		if a != expected[i] {
			t.Errorf("error in TestReadReducedDefautlts (outer qoutes inserted): expected \"%s\", got \"%s\"", expected[i], a)
		}
	}
}

const (
	examplegrubdefaults = `# If you change this file, run 'update-grub' afterwards to update
# /boot/grub/grub.cfg.
# For full documentation of the options in this file, see:
#   info -f grub -n 'Simple configuration'

GRUB_TIMEOUT_STYLE=menu
GRUB_TIMEOUT=10
GRUB_CMDLINE_LINUX_DEFAULT=""
GRUB_CMDLINE_LINUX=""
GRUB_DISABLE_SUBMENU=y

# Uncomment to enable BadRAM filtering, modify to suit your needs
# This works with Linux (no patch required) and with any kernel that obtains
# the memory map information from GRUB (GNU Mach, kernel of FreeBSD ...)
#GRUB_BADRAM="0x01234567,0xfefefefe,0x89abcdef,0xefefefef"

# Uncomment to disable graphical terminal (grub-pc only)
#GRUB_TERMINAL=console

# The resolution used on graphical terminal
# note that you can use only modes which your graphic card supports via VBE
#GRUB_GFXMODE=640x480

# Uncomment if you don't want GRUB to pass "root=UUID=xxx" parameter to Linux
GRUB_DISABLE_LINUX_UUID=true

# Uncomment to disable generation of recovery mode menu entries
#GRUB_DISABLE_RECOVERY="true"

# Uncomment to get a beep at grub start
#GRUB_INIT_TUNE="480 440 1"

GRUB_DEFAULT="gnulinux-5.11.0-advanced-b70cb823-9505-4ab6-bc0a-ca359515bf03"
`
	examplegrubdefaultsnolines = `# If you change this file, run 'update-grub' afterwards to update
# /boot/grub/grub.cfg.
# For full documentation of the options in this file, see:
#   info -f grub -n 'Simple configuration'
GRUB_TIMEOUT_STYLE=menu
GRUB_TIMEOUT=10
GRUB_CMDLINE_LINUX_DEFAULT=""
GRUB_CMDLINE_LINUX=""
GRUB_DISABLE_SUBMENU=y
# Uncomment to enable BadRAM filtering, modify to suit your needs
# This works with Linux (no patch required) and with any kernel that obtains
# the memory map information from GRUB (GNU Mach, kernel of FreeBSD ...)
#GRUB_BADRAM="0x01234567,0xfefefefe,0x89abcdef,0xefefefef"
# Uncomment to disable graphical terminal (grub-pc only)
#GRUB_TERMINAL=console
# The resolution used on graphical terminal
# note that you can use only modes which your graphic card supports via VBE
#GRUB_GFXMODE=640x480
# Uncomment if you don't want GRUB to pass "root=UUID=xxx" parameter to Linux
GRUB_DISABLE_LINUX_UUID=true
# Uncomment to disable generation of recovery mode menu entries
#GRUB_DISABLE_RECOVERY="true"
# Uncomment to get a beep at grub start
#GRUB_INIT_TUNE="480 440 1"
GRUB_DEFAULT="gnulinux-5.11.0-advanced-b70cb823-9505-4ab6-bc0a-ca359515bf03"
`

	examplegrubcfg = `### BEGIN /etc/grub.d/05_debian_theme ###
set menu_color_normal=white/black
set menu_color_highlight=black/light-gray
### END /etc/grub.d/05_debian_theme ###

### BEGIN /etc/grub.d/10_linux ###
function gfxmode {
	set gfxpayload="${1}"
	if [ "${1}" = "keep" ]; then
		set vt_handoff=vt.handoff=7
	else
		set vt_handoff=
	fi
}
if [ "${recordfail}" != 1 ]; then
  if [ -e ${prefix}/gfxblacklist.txt ]; then
    if hwmatch ${prefix}/gfxblacklist.txt 3; then
      if [ ${match} = 0 ]; then
        set linux_gfx_mode=keep
      else
        set linux_gfx_mode=text
      fi
    else
      set linux_gfx_mode=text
    fi
  else
    set linux_gfx_mode=keep
  fi
else
  set linux_gfx_mode=text
fi
export linux_gfx_mode
menuentry 'Ubuntu, with Linux 5.11.0' --class ubuntu --class gnu-linux --class gnu --class os $menuentry_id_option 'gnulinux-5.11.0-advanced-b70cb823-9505-4ab6-bc0a-ca359515bf03' {
	recordfail
	load_video
	gfxmode $linux_gfx_mode
	insmod gzio
	if [ x$grub_platform = xxen ]; then insmod xzio; insmod lzopio; fi
	insmod part_gpt
	insmod ext2
	if [ x$feature_platform_search_hint = xy ]; then
	  search --no-floppy --fs-uuid --set=root  b70cb823-9505-4ab6-bc0a-ca359515bf03
	else
	  search --no-floppy --fs-uuid --set=root b70cb823-9505-4ab6-bc0a-ca359515bf03
	fi
	echo	'Loading Linux 5.11.0 ...'
	linux	/boot/vmlinuz-5.11.0 root=/dev/nvme0n1p2 ro  
	echo	'Loading initial ramdisk ...'
	initrd	/boot/initrd.img-5.11.0
}
menuentry 'Ubuntu, with Linux 5.11.0 (recovery mode)' --class ubuntu --class gnu-linux --class gnu --class os $menuentry_id_option 'gnulinux-5.11.0-recovery-b70cb823-9505-4ab6-bc0a-ca359515bf03' {
	recordfail
	load_video
	insmod gzio
	if [ x$grub_platform = xxen ]; then insmod xzio; insmod lzopio; fi
	insmod part_gpt
	insmod ext2
	if [ x$feature_platform_search_hint = xy ]; then
	  search --no-floppy --fs-uuid --set=root  b70cb823-9505-4ab6-bc0a-ca359515bf03
	else
	  search --no-floppy --fs-uuid --set=root b70cb823-9505-4ab6-bc0a-ca359515bf03
	fi
	echo	'Loading Linux 5.11.0 ...'
	linux	/boot/vmlinuz-5.11.0 root=/dev/nvme0n1p2 ro recovery nomodeset dis_ucode_ldr 
	echo	'Loading initial ramdisk ...'
	initrd	/boot/initrd.img-5.11.0
}
menuentry 'Ubuntu, with Linux 5.11.0.old' --class ubuntu --class gnu-linux --class gnu --class os $menuentry_id_option 'gnulinux-5.11.0.old-advanced-b70cb823-9505-4ab6-bc0a-ca359515bf03' {
	recordfail
	load_video
	gfxmode $linux_gfx_mode
	insmod gzio
	if [ x$grub_platform = xxen ]; then insmod xzio; insmod lzopio; fi
	insmod part_gpt
	insmod ext2
	if [ x$feature_platform_search_hint = xy ]; then
	  search --no-floppy --fs-uuid --set=root  b70cb823-9505-4ab6-bc0a-ca359515bf03
	else
	  search --no-floppy --fs-uuid --set=root b70cb823-9505-4ab6-bc0a-ca359515bf03
	fi
	echo	'Loading Linux 5.11.0.old ...'
	linux	/boot/vmlinuz-5.11.0.old root=/dev/nvme0n1p2 ro  
	echo	'Loading initial ramdisk ...'
	initrd	/boot/initrd.img-5.11.0
}`
)
