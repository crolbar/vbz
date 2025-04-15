{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

  outputs = inputs: let
    system = "x86_64-linux";
    pkgs = import inputs.nixpkgs {inherit system;};

    deps = with pkgs; [
      miniaudio
      gcc
    ];
  in {
    devShells.${system}.default = pkgs.mkShell {
      packages = with pkgs;
        [
          clang-tools
        ]
        ++ deps;

      LD_LIBRARY_PATH = with pkgs; lib.makeLibraryPath [libpulseaudio];
    };

    packages.${system}.default = pkgs.buildGoModule {
      pname = "vbz";
      version = "0.1";
      src = ./.;
      buildInputs = with pkgs; [makeWrapper] ++ deps;
      installPhase = ''
        mkdir -p $out/bin

        dir="$GOPATH/bin"
        [ -e "$dir" ] && cp -r $dir $out

        mkdir -p $out/share/applications
        cp vbz.desktop $out/share/applications
      '';
      postFixup = ''
        wrapProgram $out/bin/vbz --set LD_LIBRARY_PATH ${pkgs.libpulseaudio}/lib
      '';
      vendorHash = "sha256-xKAAPSz1IE8VVR6sNViUicG5Wz5m0IiOkutbOPhywKQ=";
    };
  };
}
