
import subprocess


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

	print("We are running:\n  %s" % " ".join(command))
	p = subprocess.run(command, stdout=subprocess.PIPE)
	print("Done!")

	return p.stdout.decode("utf-8")


def get_stats_from_output(output):
	# separeate into lines
	lines = output.split("\n")

	# find from where overall stats are
	overall_stats = lines.index("--OVERALL STATS--")

	# empty results dict
	results = {
		"min":None,
		"max":None,
		"mean":None,
		"median":None}

	# check each line in overall stats
	for line in lines[overall_stats:]:
		# for each result we're looking for
		for key in results.keys():
			if key in line.lower():
				result_string = line.split(":")[-1].strip()
				results[key] = float(result_string.replace("ms",""))

	return results



if __name__ == "__main__":
	output = run_crb("http://localhost", count=10)
	results = get_stats_from_output(output)

	print(results)