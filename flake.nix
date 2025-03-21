{
  description = "gowindows development";

  inputs = {
    nixpkgs.url = "nixpkgs";
    pre-commit-hooks.url = "github:cachix/git-hooks.nix";
  };

  outputs = {
    self,
    nixpkgs,
    ...
  } @ inputs: let
    allSystems = [
      "x86_64-linux"
      "aarch64-linux"
      "x86_64-darwin"
      "aarch64-darwin"
    ];

    forAllSystems = f:
      nixpkgs.lib.genAttrs allSystems (system:
        f {
          pkgs = import nixpkgs {inherit system;};
        });
  in {
    checks = forAllSystems (system: {
      pre-commit-check = inputs.pre-commit-hooks.lib.${system}.run {
        src = ./.;
        hooks = {
          nixpkgs-fmt.enable = true;
        };
      };
    });
    devShells = forAllSystems ({pkgs}: {
      default = pkgs.mkShell {
        packages = with pkgs; [
          go
          gotools
          golangci-lint
          pre-commit
        ];

        shellHook = let
          pre-commit = "${pkgs.pre-commit}/bin/pre-commit";
        in
          /*
          bash
          */
          ''
            if [ ! -e .git/hooks/pre-commit ]; then
              echo "Setting up pre-commit hooks..."
               ${pre-commit} install
            else
              echo "Pre-commit hooks already installed."
            fi
          '';
      };
    });
  };
}
