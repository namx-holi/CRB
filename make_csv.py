
import subprocess
import time


def estimate_completion_times(counts, loops, cooldowns, resp_times):
	print("--Estimated completion times--")
	for resp_time in resp_times:
		estimated_time = loops * (cooldown + resp_time*1000) * len(counts)

		print("  %4dms average response time will take %d seconds (%.2f minutes)" % (
				resp_time*1000, estimated_time/1000, estimated_time/60000))
	print("\n")


def run_full_benchmark(url, counts, looking_for, loops=5, cooldown=3000):


	estimate_completion_times(counts, loops, cooldown, [0, 1, 2, 3])

	results = {}

	for count in counts:
		output = run_crb(url, loops=loops, count=count, cooldown=cooldown)

		# sleep for a cooldown. this is because there's no cooldown time
		# at the end of a run of crb
		time.sleep(cooldown/1000.0)

		results[count] = get_stats_from_output(output, looking_for)

	print("Done!")

	return results


def run_crb(url, count=1, loops=1, cooldown=1000, verbose=False):
	command = ["crb"]

	if url:
		command.append("-url=%s" % url)
	
	if count:
		command.append("-count=%d" % count)
	
	if loops:
		command.append("-loops=%d" % loops)
	
	if cooldown:
		command.append("-cooldown=%d" % cooldown)
	
	if verbose:
		command.append("-verbose")

	print("We are running:\n  %s\n" % " ".join(command))
	p = subprocess.run(command, stdout=subprocess.PIPE)

	return p.stdout.decode("utf-8")


def get_stats_from_output(output, looking_for):
	# separeate into lines
	lines = output.split("\n")

	# find from where overall stats are
	overall_stats = lines.index("--OVERALL STATS--")

	# empty results dict
	results = {}
	for key in looking_for:
		results[key] = None

	# check each line in overall stats
	for line in lines[overall_stats:]:
		# for each result we're looking for
		for key in results.keys():
			if key.lower() in line.lower():
				result_string = line.split(":")[-1].strip()
				results[key] = float(result_string.replace("ms",""))

	return results


def stats_to_csv(stats, counts, looking_for):
	# counts and looking for are added
	# so data can be in order

	delim = ","

	rows = []

	# header
	header_items = ["NbResponse"] + [key.title() for key in looking_for]
	rows.append(delim.join(header_items))

	for count in counts:
		line_items = [count]

		for key in looking_for:
			line_items.append(stats[count][key])

		rows.append(delim.join(
			str(line_item) for line_item in line_items))

	return "\n".join(rows)


def write_csv_to_file(csv_text, filename):
	f = open(filename, "w")
	f.write(csv_text)
	f.close()


if __name__ == "__main__":
	counts = [1, 2, 3, 4, 5, 10, 15, 20, 30, 40, 50, 75, 100, 150, 200]
	looking_for = ["Min","Max","Mean","Median"]
	url = "http://localhost"
	loops = 5
	cooldown = 500

	results = run_full_benchmark(url, counts, looking_for,
		loops=loops, cooldown=cooldown)

	print("\n\n\n")

	csv_text = stats_to_csv(results, counts, looking_for)
	print(csv_text)
	write_csv_to_file(csv_text, "results.csv")
