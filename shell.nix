{ sources ? import ./nix/sources.nix, pkgs ? import ./nix { inherit sources; } }:

pkgs.mkShell {
  name = "steamcmd-shell";

  buildInputs = with pkgs; [
    go
  ];
}
