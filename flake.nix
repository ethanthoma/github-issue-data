{
  description = "A basic gomod2nix flake";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs = {
        nixpkgs.follows = "nixpkgs";
        flake-utils.follows = "flake-utils";
      };
    };
  };

  outputs = { self, nixpkgs, flake-utils, gomod2nix }:
    (flake-utils.lib.eachDefaultSystem
      (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};

          callPackage = pkgs.darwin.apple_sdk_11_0.callPackage or pkgs.callPackage;

          packageNames = [ "repos" "comments" "history" "sample" "stargazers" ];

          buildGoPackage = name: (
            callPackage ./nix/template.nix {
              inherit name pkgs;
              inherit (gomod2nix.legacyPackages.${system}) buildGoApplication;
            }
          );
        in
        {
          packages = builtins.listToAttrs (builtins.map
            (name: { name = name; value = buildGoPackage name; })
            packageNames
          );

          devShells.default = callPackage ./nix/shell.nix {
            inherit (gomod2nix.legacyPackages.${system}) mkGoEnv gomod2nix;
          };
        }
      )
    );
}
