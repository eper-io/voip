module gitlab.com/eper.io/engine

// Licensed under Creative Commons CC0.
//
// To the extent possible under law, the author(s) have dedicated all copyright
// neighboring rights to this software to the public domain worldwide.
// This software is distributed without any warranty.
// You should have received a copy of the CC0 Public Domain Dedication along wi
// If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

// The project is designed to use 'go mod vendor'
// Do not use 'go mod tidy' or 'go mod download' for security compliance
// We also avoid using modules altogether to minimize the risk of third party vulnerabilites.
// The codebase is few thousand lines so that it can be security reviewed easily.
// TODO Do we want to stick to the latest 1.19?

go 1.19
