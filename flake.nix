{
  description = "gowindows development";

  inputs.nixpkgs.url = "nixpkgs";

  outputs = {
    self,
    nixpkgs,
    ...
  }: let
    systems = ["x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin"];
  in {
    devShells = nixpkgs.lib.genAttrs systems (system: let
      pkgs = import nixpkgs {inherit system;};
    in {
      default = pkgs.mkShell {
        packages = with pkgs; [
          go
          gotools
          golangci-lint
          pre-commit
          nodejs
          nodePackages.npm
        ];

        shellHook = ''
          HOOK_PATH=$(git rev-parse --git-path hooks/pre-commit)
          if [ ! -f "$HOOK_PATH" ]; then
            echo "Setting up pre-commit hooks..."
            ${pkgs.pre-commit}/bin/pre-commit install
          fi
        '';
      };
    });
  };
}
