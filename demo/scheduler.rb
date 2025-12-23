# frozen_string_literal: true

require 'sidekiq/api'

class JobScheduler
  QUEUES = %w[critical default mailers batch low].freeze
  MAX_JOBS_PER_QUEUE = 10_000
  MAX_RETRY_QUEUE = 20_000
  SCHEDULE_BATCH_SIZE = 100

  JOB_DEFINITIONS = [
    { job: 'EmailDeliveryJob', queue: 'mailers', weight: 15,
      args: -> { [fake_email, fake_subject, %w[welcome reset_password invoice].sample] } },
    { job: 'PaymentProcessingJob', queue: 'critical', weight: 5,
      args: -> { [rand(100_000..999_999), rand(10.0..500.0).round(2), %w[USD EUR GBP].sample] } },
    { job: 'ImageProcessingJob', queue: 'default', weight: 10,
      args: -> { [rand(1..100_000), %w[resize crop thumbnail watermark].sample(rand(1..3))] } },
    { job: 'ReportGenerationJob', queue: 'batch', weight: 3,
      args: -> { [%w[sales inventory users activity].sample, "#{rand(1..12)}/2024", rand(1..1000)] } },
    { job: 'NotificationJob', queue: 'critical', weight: 20,
      args: -> { [rand(1..100_000), %w[push sms in_app].sample, { "message" => fake_message }] } },
    { job: 'DataSyncJob', queue: 'batch', weight: 5,
      args: -> { [%w[salesforce hubspot stripe].sample, 'local_db', Array.new(rand(1..10)) { rand(1..10_000) }] } },
    { job: 'CacheWarmupJob', queue: 'low', weight: 8,
      args: -> { ["cache:#{%w[products users orders].sample}:#{rand(1..1000)}", %w[product user order].sample] } },
    { job: 'AnalyticsJob', queue: 'low', weight: 25,
      args: -> { [%w[page_view click purchase signup].sample, { "page" => "/page/#{rand(1..100)}" }, Time.now.to_i] } },
    { job: 'CleanupJob', queue: 'low', weight: 4,
      args: -> { [%w[temp_files sessions logs exports].sample, rand(7..90)] } },
    { job: 'WebhookDeliveryJob', queue: 'default', weight: 10,
      args: -> { ["https://example.com/webhooks/#{rand(1..100)}", %w[order.created user.updated payment.completed].sample, { "id" => rand(1..10_000) }] } }
  ].freeze

  def self.fake_email
    "user#{rand(1..100_000)}@example.com"
  end

  def self.fake_subject
    ['Welcome!', 'Your order confirmation', 'Password reset', 'Weekly digest', 'Invoice #%d' % rand(1000..9999)].sample
  end

  def self.fake_message
    ['New message received', 'Your order shipped', 'Payment confirmed', 'Reminder: complete your profile'].sample
  end

  def initialize
    @weighted_jobs = build_weighted_job_list
    @running = false
  end

  def start
    @running = true
    puts "Starting job scheduler..."
    puts "Max jobs per queue: #{MAX_JOBS_PER_QUEUE}"
    puts "Max retry queue: #{MAX_RETRY_QUEUE}"
    puts "Queues: #{QUEUES.join(', ')}"

    while @running
      schedule_jobs_if_needed
      sleep 0.5
    end
  end

  def stop
    @running = false
  end

  private

  def build_weighted_job_list
    JOB_DEFINITIONS.flat_map { |defn| Array.new(defn[:weight], defn) }
  end

  def schedule_jobs_if_needed
    retry_size = Sidekiq::RetrySet.new.size

    # Pause scheduling if retry queue is too large
    if retry_size >= MAX_RETRY_QUEUE
      puts "[#{Time.now.strftime('%H:%M:%S')}] Retry queue full (#{retry_size}), pausing scheduling..."
      return
    end

    queue_sizes = fetch_queue_sizes

    QUEUES.each do |queue_name|
      current_size = queue_sizes[queue_name] || 0
      available_capacity = MAX_JOBS_PER_QUEUE - current_size

      next if available_capacity <= 0

      jobs_to_schedule = [available_capacity, SCHEDULE_BATCH_SIZE].min
      schedule_jobs_for_queue(queue_name, jobs_to_schedule)
    end
  end

  def fetch_queue_sizes
    QUEUES.each_with_object({}) do |queue_name, sizes|
      sizes[queue_name] = Sidekiq::Queue.new(queue_name).size
    end
  end

  def schedule_jobs_for_queue(queue_name, count)
    queue_jobs = @weighted_jobs.select { |j| j[:queue] == queue_name }
    return if queue_jobs.empty?

    count.times do
      job_def = queue_jobs.sample
      job_class = Object.const_get(job_def[:job])
      job_class.set(cattr: {tenanant_id: rand(1..10000)}).perform_async(*job_def[:args].call)
    end

    puts "[#{Time.now.strftime('%H:%M:%S')}] Scheduled #{count} jobs for queue '#{queue_name}'"
  end
end
