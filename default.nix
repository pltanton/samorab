with import <nixpkgs> {};
buildGoPackage {
  name = "samorab";
  version = "1.0";
  goPackagePath = "github.com/pltanton/samorab";
  src = lib.cleanSourceWith {
    filter = (path: type:
      ! (builtins.any
          (r: (builtins.match r (builtins.baseNameOf path)) != null)
          [
            "\.env"
            ".go"
          ])
    );
    src = lib.cleanSource ./.;
  };
  goDeps = ./deps.nix;
}
