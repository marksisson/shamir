{ self, ... }: {
  perSystem = { inputs', ... }:
    let
      inherit (inputs'.gomod2nix.legacyPackages) buildGoApplication;
    in
    {
      packages.shamir = buildGoApplication {
        pname = "shamir";
        version = "0.0.1";
        src = self;
        modules = "${self}/gomod2nix.toml";
      };
    };
}
