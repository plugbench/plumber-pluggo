{
  description = "TODO: fill me in";
  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = { self, nixpkgs, flake-utils }:
    (flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        nats-plumber = pkgs.callPackage ./derivation.nix {};
      in {
        packages = {
          default = nats-plumber;
          inherit nats-plumber;
        };
        checks = {
          test = pkgs.runCommandNoCC "nats-plumber-test" {} ''
            mkdir -p $out
            : ${nats-plumber}
          '';
        };
    })) // {
      overlays.default = final: prev: {
        nats-plumber = prev.callPackage ./derivation.nix {};
      };
    };
}
