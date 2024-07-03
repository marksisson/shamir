{ inputs, ... }: {
  imports = [
    inputs.treefmt-nix.flakeModule
  ];

  perSystem = { pkgs, ... }: {
    treefmt = {
      projectRootFile = "flake.nix";
      package = pkgs.treefmt;
      programs = {
        nixpkgs-fmt.enable = true;
        shfmt.enable = true;
      };
    };
  };
}
