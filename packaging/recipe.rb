class Vell < FPM::Cookery::Recipe
  GOPACKAGE = 'vell'

  name         'vell'
  version      '1.1.0'
  source       'https://github.com/rkcpi/vell.git', :with => :git, :tag => @version
  revision     '1'
  description  'Lightweight repository management tool for RPM repositories'
  homepage     'https://rkcpi.github.io/vell/'

  build_depends %w(golang git)
  depends       'createrepo'

  post_install 'vell.postinst'

  def build
    pkgdir = builddir("gobuild/src/#{GOPACKAGE}")
    mkdir_p pkgdir
    cp_r Dir["*"], pkgdir

    ENV["GOPATH"] = [
      builddir("gobuild"),
    ].join(":")

    safesystem "go version"
    safesystem "go env"
    safesystem "go get -v #{GOPACKAGE}/..."
  end

  def install
    bin.install builddir("gobuild/bin/vell")
    root('lib/systemd/system').install_p workdir('vell.service')
  end

end
