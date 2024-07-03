{ self, ... }: {
  perSystem = { config, inputs', pkgs, ... }:
    let
      inherit (inputs'.gomod2nix.legacyPackages) mkGoEnv gomod2nix;
    in
    {
      devShells.developer =
        let
          goEnv = mkGoEnv {
            pwd = self;
            modules = "${self}/gomod2nix.toml";
          };
          formatters = [ treefmt ] ++ treefmt-programs;
          tools = [ goEnv gomod2nix ];
          treefmt = config.treefmt.build.wrapper;
          treefmt-programs = builtins.attrValues config.treefmt.build.programs;
        in
        pkgs.mkShell {
          packages = formatters ++ tools;
          shellHook = ''
            ${config.pre-commit.installationScript}
          '';
        };
    };
}
