{ stdenv, lib, ... }:

stdenv.mkDerivation rec {
  pname = "nats-plumber";
  version = "0.1.0";

  src = ./.;

  meta = with lib; {
    description = "TODO: fill me in";
    homepage = "https://github.com/eraserhd/nats-plumber";
    license = licenses.publicDomain;
    platforms = platforms.all;
    maintainers = [ maintainers.eraserhd ];
  };
}
