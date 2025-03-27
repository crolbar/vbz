{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

  outputs = inputs: let
    system = "x86_64-linux";
    pkgs = import inputs.nixpkgs {inherit system;};
  in {
    devShells.${system}.default = pkgs.mkShell {
      packages = with pkgs; [
        miniaudio
        gcc
        clang-tools
      ];

      shellHook = ''
        export MINIAUDIO="${pkgs.miniaudio}"
      '';

      LD_LIBRARY_PATH = with pkgs; lib.makeLibraryPath [alsa-lib libpulseaudio];
    };
  };
}
