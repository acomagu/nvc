{
  description = "nvc";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.05";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "nvc";
          version = "0.1.0";

          src = self;

          vendorHash = "sha256-sr1GVWu637dHKTG6KCRPKygxpr0O5Ckfk+3LaXJb9tw=";
        };

        apps.default = {
          type = "app";
          program = "${self.packages.${system}.default}/bin/nvc";
        };
      });
}

