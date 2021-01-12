require 'fileutils'

def kill_process_by_name name
  out = `ps aux | grep #{name}`
  out.each_line do |line|
    id = line.split()[1]
    puts "Killing #{id} for relating to #{name}"
    `kill '#{id}'`
  end
end

if __FILE__ == $0
  execute_time = Time.new.strftime("%Y-%m-%d-%H:%M:%S:%L")
  Dir.mkdir("log") unless Dir.exist?"log" #DEBUG
  # Dir.mkdir("log/#{execute_time}")
  Dir.mkdir("log/current") unless Dir.exist? "log/current"
  # `mkdir log/#{execute_time}`

  # clean up record
  record_dir = File.join(__dir__, 'record')
  old_record_dir = File.join(record_dir, 'old_record')
  FileUtils.mkdir_p(old_record_dir) unless File.directory?(old_record_dir)
  Dir.new(record_dir).each {|file|
    file = File.join(record_dir, file)
    FileUtils.mv(file, old_record_dir) if File.file? file
  }

  # start env
  #env_id = spawn "go run environment/main.go > log/#{execute_time}/environment.log 2>&1"
  env_id = spawn "go run environment/main.go > log/current/environment.log 2>&1" #DEBUG change to one folder for easy of debugging

  # wait a bit
  sleep 5

  # start agents
  agents_id = []
  for i in (1..ARGV[0].to_i)
    #agents_id.append(spawn("python3 agent/main.py > log/#{execute_time}/agent#{i}.log 2>&1"))
    agents_id.append(spawn("python3 agent/main.py > log/current/agent#{i}.log 2>&1")) #DEBUG
  end

  # at_exit{
  #   agents_id.each do |id|
  #   end
  #   kill_process_by_name "main.py"
  #   kill_process_by_name "main.go"
  # }

  while true
  end

end
