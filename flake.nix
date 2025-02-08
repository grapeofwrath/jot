{
  description = "A Zettelkasten CLI, written in Go";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = {
    self,
    nixpkgs,
  }: let
    allSystems = [
      "x86_64-linux"
      "aarch64-linux"
      "x86_64-darwin"
      "aarch64-darwin"
    ];
    forAllSystems = f:
      nixpkgs.lib.genAttrs allSystems (system:
        f {
          pkgs = import nixpkgs {inherit system;};
        });
  in {
    packages = forAllSystems ({pkgs}: {
      default = pkgs.buildGoModule rec {
        pname = "jot";
        version = "1.0.0";

        src = ./.;

        vendorHash = null;

        meta = {
          description = "A Zettelkasten CLI, written in Go";
          homepage = "https:github.com/grapeofwrath/jot";
          license = pkgs.lib.licenses.gpl3Plus;
        };
      };
    });
  };
}
