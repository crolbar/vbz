{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

  outputs = inputs: let
    systems = ["x86_64-linux" "aarch64-linux"];
    forEachSystem = inputs.nixpkgs.lib.genAttrs systems;
    pkgsFor = inputs.nixpkgs.legacyPackages;
  in {
    devShells = forEachSystem (system: let
      pkgs = pkgsFor.${system};
    in {
      default = pkgs.mkShell {
        packages = with pkgs; [
          clang-tools
          miniaudio
        ];

        LD_LIBRARY_PATH = with pkgs; lib.makeLibraryPath [libpulseaudio];
      };
    });

    packages = forEachSystem (system: let
      pkgs = pkgsFor.${system};
    in {
      default = pkgs.buildGoModule {
        pname = "vbz";
        version = "0.1";
        src = ./.;
        buildInputs = with pkgs; [makeWrapper miniaudio];
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
    });
  };
}
