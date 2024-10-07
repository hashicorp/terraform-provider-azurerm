schema_version = 1

project {
  license        = "MPL-2.0"
  copyright_year = 2023

  header_ignore = [
    ".ci/**",
    ".github/**",
    ".teamcity/**",
    ".release/**",
    "vendor/**",
    "internal/**/**_gen.go", # Pandora generated files
    ".goreleaser.yml",
  ]
}
