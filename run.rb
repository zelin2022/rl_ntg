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
  # Dir.mkdir("log") unless Dir.exist?"log"
  # Dir.mkdir("log/#{execute_time}")
  `mkdir log/#{execute_time}`

  # start env
  env_id = spawn "go run environment/main.go > log/#{execute_time}/environment.log 2>&1"

  # wait a bit
  sleep 5

  # start agents
  agents_id = []
  for i in (1..ARGV[0].to_i)
    agents_id.append(spawn("python3 agent/main.py > log/#{execute_time}/agent#{i}.log 2>&1"))
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
