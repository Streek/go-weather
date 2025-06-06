.PHONY: build clean install package release version-bump

# Build variables
BINARY_NAME=wgo
VERSION=$(shell cat VERSION)
BUILD_DIR=build
LDFLAGS=-X 'main.appVersion=$(VERSION)'

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) main.go

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

# Version bump helpers
version-bump-patch:
	@current=$$(cat VERSION); \
	new=$$(echo $$current | awk -F. '{$$3++; print $$1"."$$2"."$$3}'); \
	echo $$new > VERSION; \
	echo "Version bumped from $$current to $$new"

version-bump-minor:
	@current=$$(cat VERSION); \
	new=$$(echo $$current | awk -F. '{$$2++; $$3=0; print $$1"."$$2"."$$3}'); \
	echo $$new > VERSION; \
	echo "Version bumped from $$current to $$new"

version-bump-major:
	@current=$$(cat VERSION); \
	new=$$(echo $$current | awk -F. '{$$1++; $$2=0; $$3=0; print $$1"."$$2"."$$3}'); \
	echo $$new > VERSION; \
	echo "Version bumped from $$current to $$new"
