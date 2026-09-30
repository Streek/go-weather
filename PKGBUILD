# Maintainer: Keith Connolly <streek@mesaro.com>
pkgname=go-weather
pkgver=1.0.2
pkgrel=1
pkgdesc="Console-based weather application"
arch=('x86_64' 'aarch64')
url="https://github.com/streek/go-weather"
license=('GPL3')
depends=()
makedepends=('go>=1.21')
source=("$pkgname-$pkgver.tar.gz::$url/archive/v$pkgver.tar.gz")
sha256sums=('e0d070fc77058b3f1b7bb713370573267a710277c401b8e867a523291f3aa6dc')

build() {
  cd "$srcdir/go-weather-$pkgver"
  export CGO_CPPFLAGS="${CPPFLAGS}"
  export CGO_CFLAGS="${CFLAGS}"
  export CGO_CXXFLAGS="${CXXFLAGS}"
  export CGO_LDFLAGS="${LDFLAGS}"
  export GOFLAGS="-buildmode=pie -trimpath -mod=readonly -modcacherw"
  
  # Read version and inject it during build
  VERSION=$(cat VERSION)
  go build -ldflags "-X 'main.appVersion=$VERSION'" -o go-weather .
}

package() {
  cd "$srcdir/go-weather-$pkgver"
  install -Dm755 go-weather "$pkgdir/usr/bin/go-weather"
  install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE"
}
