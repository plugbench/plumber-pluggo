{
  description = "Re-implementation of the Plan9 plumber using NATS.io";
  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = { self, nixpkgs, flake-utils }:
    (flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        plumber-pluggo = pkgs.callPackage ./derivation.nix {};
      in {
        packages = {
          default = plumber-pluggo;
          inherit plumber-pluggo;
        };
        checks = {
          test = pkgs.runCommand "plumber-pluggo-test" {} ''
            mkdir -p $out
            : ${plumber-pluggo}
          '';
        };
        devShells.default = pkgs.mkShell {
          nativeBuildInputs = with pkgs; [
            go_1_19
          ];
        };
    })) // {
      overlays.default = final: prev: {
        plumber-pluggo = prev.callPackage ./derivation.nix {};
      };
    };
}
