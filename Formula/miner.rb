# Homebrew Formula for Miner
# Installation: brew install YOUR_USERNAME/tap/miner
# 
# To create your tap:
# 1. Create a GitHub repo: YOUR_USERNAME/homebrew-tap
# 2. Add this file as: Formula/miner.rb
# 3. Users install with: brew tap YOUR_USERNAME/tap && brew install miner

class Miner < Formula
  desc "Standalone database manager with Adminer and FrankenPHP"
  homepage "https://github.com/YOUR_USERNAME/miner"
  version "1.0.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/YOUR_USERNAME/miner/releases/download/v1.0.0/miner-darwin-arm64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256_ARM64"
    else
      url "https://github.com/YOUR_USERNAME/miner/releases/download/v1.0.0/miner-darwin-amd64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256_AMD64"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/YOUR_USERNAME/miner/releases/download/v1.0.0/miner-linux-arm64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256_LINUX_ARM64"
    else
      url "https://github.com/YOUR_USERNAME/miner/releases/download/v1.0.0/miner-linux-amd64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256_LINUX_AMD64"
    end
  end

  depends_on "frankenphp" => :recommended

  def install
    bin.install "miner-#{OS.kernel_name.downcase}-#{Hardware::CPU.arch}" => "miner"
  end

  def post_install
    ohai "Miner installed successfully!"
    ohai "Run 'sudo miner install' to set up hosts file and CLI commands"
  end

  def caveats
    <<~EOS
      To complete installation:
        1. sudo miner install
        2. miner
        3. Open http://miner.local:88

      Note: FrankenPHP is recommended but will be auto-installed if missing.
    EOS
  end

  test do
    assert_match "Miner v", shell_output("#{bin}/miner version")
  end
end
