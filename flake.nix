{
  description = "tcpeek - TCP event listener daemon";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachSystem [ "x86_64-linux" "aarch64-linux" ] (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "tcpeek";
          version = "0.1.0";
          src = ./.;
          vendorHash = "sha256-6Uin5Zla/AbtKKZRdV26TkNk7ITxxIUw8jUkZLQjnjg=";
          subPackages = [ "cmd/tcpeek" ];
        };

        devShells.default = pkgs.mkShell {
          buildInputs = [ pkgs.go ];
        };
      });
}
