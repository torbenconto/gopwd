package main

import "github.com/blang/semver/v4"

func GetVersion() semver.Version {
	version := semver.Version{
		Major: 1,
		Minor: 1,
		Patch: 9,
		Pre: []semver.PRVersion{
			{VersionStr: "git"},
		},
		Build: []string{"HEAD"},
	}
	return version
}
