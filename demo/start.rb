#!/usr/bin/env ruby
# frozen_string_literal: true

require_relative 'boot'
require_relative 'scheduler'

# Determine mode
mode = ARGV[0] || 'all'

case mode
when 'worker'
  # Start Sidekiq worker only
  puts "Starting Sidekiq worker..."
  exec('bundle', 'exec', 'sidekiq', '-r', './boot.rb', '-C', 'config/sidekiq.yml')

when 'scheduler'
  # Start scheduler only
  scheduler = JobScheduler.new
  trap('INT') { scheduler.stop }
  trap('TERM') { scheduler.stop }
  scheduler.start

when 'web'
  # Start Sidekiq Web UI only
  puts "Starting Sidekiq Web UI on http://localhost:9292..."
  exec('bundle', 'exec', 'rackup', '-o', '0.0.0.0', '-p', '9292')

when 'all'
  # Start worker, scheduler, and web UI
  puts "Starting Sidekiq simulation..."
  puts "Sidekiq Web UI: http://localhost:9292"
  puts "Press Ctrl+C to stop"
  puts ""

  pids = []

  # Fork worker process
  pids << fork do
    exec('bundle', 'exec', 'sidekiq', '-r', './boot.rb', '-C', 'config/sidekiq.yml')
  end

  # Fork web UI process
  pids << fork do
    exec('bundle', 'exec', 'rackup', '-o', '0.0.0.0', '-p', '9292')
  end

  # Give processes time to start
  sleep 2

  # Run scheduler in main process
  scheduler = JobScheduler.new

  shutdown = proc do
    puts "\nShutting down..."
    scheduler.stop
    pids.each { |pid| Process.kill('TERM', pid) rescue nil }
  end

  trap('INT', &shutdown)
  trap('TERM', &shutdown)

  scheduler.start

  # Wait for child processes
  pids.each { |pid| Process.wait(pid) rescue nil }

else
  puts "Usage: ruby start.rb [mode]"
  puts "Modes:"
  puts "  all       - Start worker, scheduler, and web UI (default)"
  puts "  worker    - Start Sidekiq worker only"
  puts "  scheduler - Start job scheduler only"
  puts "  web       - Start Sidekiq Web UI only"
  exit 1
end
