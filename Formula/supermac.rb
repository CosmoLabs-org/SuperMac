class Supermac < Formula
  desc "macOS system management CLI — 12 modules, 127 commands"
  homepage "https://github.com/CosmoLabs-org/SuperMac"
  url "https://github.com/CosmoLabs-org/SuperMac/archive/refs/tags/v0.2.0.tar.gz"
  sha256 "PLACEHOLDER_UPDATE_ON_RELEASE"
  license "MIT"
  head "https://github.com/CosmoLabs-org/SuperMac.git", branch: "master"

  depends_on "go" => :build
  depends_on :macos

  def install
    cd "supermac-go" do
      system "go", "build",
        "-ldflags",
        "-X github.com/cosmolabs-org/supermac/internal/version.Version=#{version}",
        "-o", bin/"mac",
        "./cmd/mac"
    end
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/mac version")
  end
end
