{
  imports = [
    ./shamir
  ];

  perSystem = { self', ... }: {
    packages.default = self'.packages.shamir;
  };
}
