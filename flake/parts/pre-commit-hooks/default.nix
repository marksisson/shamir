{ inputs, ... }: {
  imports = [
    inputs.git-hooks.flakeModule
  ];

  perSystem = { config, ... }: {
    pre-commit.settings.hooks = {
      treefmt.enable = true;
      treefmt.package = config.treefmt.build.wrapper;
    };
  };
}
