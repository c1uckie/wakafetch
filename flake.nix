{
  description = "Wakafetch build with nix";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = import nixpkgs { inherit system; };

      in {
        packages.default = pkgs.buildGoModule {
          pname = "wakafetch";
          name = "wakafetch";

          src = pkgs.fetchFromGitHub {
            owner = "marc55s";
            repo = "wakafetch";
            rev = "d182dcef64945033ea256e0ccfd8bdb36382bc4c";
            sha256 ="sha256-smmmN1jDJpaZPXEQw+T4t2AJtMoNCaJoMFYYBuqz1Mo="; 
          };
          vendorHash = null;
        };
      });
}
