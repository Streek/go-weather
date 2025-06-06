.PHONY: build clean install package release

# Build variables
BINARY_NAME=wgo
VERSION=1.0.1
BUILD_DIR=build

build:
	go build -o $(BINARY_NAME) main.go

install: build
	install -Dm755 $(BINARY_NAME) $(DESTDIR)/usr/bin/$(BINARY_NAME)
	install -Dm644 README.md $(DESTDIR)/usr/share/doc/go-weather/README.md
	install -Dm644 LICENSE $(DESTDIR)/usr/share/licenses/go-weather/LICENSE

clean:
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)
	rm -f *.tar.gz
	rm -f *.pkg.tar.zst

package:
	mkdir -p $(BUILD_DIR)
	makepkg -f

release:
	mkdir -p $(BUILD_DIR)
	git archive --prefix=go-weather-$(VERSION)/ -o $(BUILD_DIR)/go-weather-$(VERSION).tar.gz HEAD
	gpg --detach-sign --armor $(BUILD_DIR)/go-weather-$(VERSION).tar.gz
