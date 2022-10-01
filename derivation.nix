{ buildGo119Module
, fetchFromGitHub
}:

buildGo119Module {
  pname = "plumber-pluggo";
  version = "0.1.0";
  src = ./.;
  vendorSha256 = "a2nddJIM5Ui1maLz2hZUssksXmvI+VQyqc5aMfZ1TbY=";
}

