{
  description = "Re-implementation of the Plan9 plumber using NATS.io";
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
        devShells.default = pkgs.mkShell {
          nativeBuildInputs = with pkgs; [
            go_1_19
          ];
        };
    })) // {
      overlays.default = final: prev: {
        nats-plumber = prev.callPackage ./derivation.nix {};
      };
    };
}
