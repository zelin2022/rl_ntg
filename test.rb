#!/usr/bin/env ruby

environment_dir =  File.join(__dir__, 'environment')

system("cd #{environment_dir} && go test ./...")

puts "GO tests passed" if $?.exitstatus == 0
puts "GO tests failed" if $?.exitstatus != 0
exit $?.exitstatus
