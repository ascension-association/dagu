all: _gokrazy/extrafiles_arm64.tar _gokrazy/extrafiles_amd64.tar

_gokrazy/extrafiles_amd64.tar:
	mkdir -p _gokrazy/extrafiles_amd64/usr/local/bin
	curl -fsSL https://github.com/dagu-org/dagu/releases/download/v1.23.4/dagu_1.23.4_linux_amd64.tar.gz | tar xzv -C _gokrazy/extrafiles_amd64/usr/local/bin/ dagu
	cd _gokrazy/extrafiles_amd64 && tar cf ../extrafiles_amd64.tar *
	rm -rf _gokrazy/extrafiles_amd64

_gokrazy/extrafiles_arm64.tar:
	mkdir -p _gokrazy/extrafiles_arm64/usr/local/bin
	curl -fsSL https://github.com/dagu-org/dagu/releases/download/v1.23.4/dagu_1.23.4_linux_arm64.tar.gz | tar xzv -C _gokrazy/extrafiles_arm64/usr/local/bin/ dagu
	cd _gokrazy/extrafiles_arm64 && tar cf ../extrafiles_arm64.tar *
	rm -rf _gokrazy/extrafiles_arm64

clean:
	rm -f _gokrazy/extrafiles_*.tar