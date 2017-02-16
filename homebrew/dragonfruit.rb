require "language/go"

class Dragonfruit < Formula
  homepage "https://github.com/ideo/dragonfruit"

  url "https://github.com/ideo/dragonfruit-cli.git", :using => :git
  version "0.5.2"

  # bottle do
  #  cellar :any
  #  sha1 "3eccd8879767b3fc95048f7d40a60edb69741892" => :yosemite
  #  sha1 "9e65954b48b3feb53ba14aafcd5081a5ede54809" => :mavericks
  #  sha1 "5f22268864f33697d6988a49c4c73ddc7368cfa2" => :mountain_lion
  # end

  #option "without-completions", "Disable bash/zsh completions"

  depends_on "go" => :build

  depends_on :hg => :build
  depends_on "couchdb"

  def install
    # start couchdb
    system "couchdb", "-b"

    system "sleep", "3s"

    # couch setup
    system "curl", "-X", "PUT", "http://localhost:5984/swagger_docs"



    # set up gopath for local installation
    ENV["GOPATH"] = buildpath
    ENV["PATH"] = ENV["PATH"] + ":" + buildpath + "/bin"

    # install godep

    # retrieve Dragonfruit CLI and its dependencies
    mkdir_p buildpath/"src/github.com/ideo/"
    ln_sf buildpath, buildpath/"src/github.com/ideo/dragonfruit-cli"
    Language::Go.stage_deps resources, buildpath/"src"
    # system "cp",  "bin/dragonfruit-cli", "dragonfruit"

    system "go", "build", "-o", "dragonfruit", "."

    etc.install "etc/dragonfruit.conf"

    # Build and install dragonfruit
    bin.install buildpath + "dragonfruit"
    
  end

  test do
    system "#{bin}/dragonfruit", "-add=false", "-serve=false"
  end
end
