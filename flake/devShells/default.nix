{
  imports = [
    ./developer
  ];

  perSystem = { self', ... }: {
    devShells.default = self'.devShells.developer;
  };
}
