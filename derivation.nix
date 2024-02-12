{ buildGo121Module
, fetchFromGitHub
}:

buildGo121Module {
  pname = "plumber-pluggo";
  version = "0.1.0";
  src = ./.;
  vendorHash = "sha256:a2nddJIM5Ui1maLz2hZUssksXmvI+VQyqc5aMfZ1TbY=";
}

