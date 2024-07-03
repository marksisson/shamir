{
  description = "";

  outputs = inputs:
    inputs.flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [ ./. { debug = true; } ];
      systems = import inputs.systems;
    };

  inputs = {
    # nix tools
    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs = {
        nixpkgs-lib.follows = "nixpkgs";
      };
    };
    git-hooks = {
      url = "github:cachix/git-hooks.nix";
      inputs = {
        nixpkgs.follows = "nixpkgs";
      };
    };
    gomod2nix = {
      #url = "github:nix-community/gomod2nix";
      url = "github:marksisson/gomod2nix?ref=patches/fix-mkgoenv-to-allow-passing-more-attributes";
      inputs = {
        nixpkgs.follows = "nixpkgs";
      };
    };
    nixpkgs = {
      url = "https://flakehub.com/f/NixOS/nixpkgs/*.tar.gz";
    };
    systems = {
      url = "github:nix-systems/default";
    };
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs = {
        nixpkgs.follows = "nixpkgs";
      };
    };
  };

  nixConfig = {
    extra-substituters = [
      "https://nix-community.cachix.org"
    ];
    extra-trusted-public-keys = [
      "nix-community.cachix.org-1:mB9FSh9qf2dCimDSUo8Zy7bkq5CX+/rkCWyvRCYg3Fs="
    ];
  };
}
