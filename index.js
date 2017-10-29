const wilson = (upvotes, downvotes, z = 1.644853) => {
  // Use 1.96 for a confidence level of 0.95.
  const n = upvotes + downvotes
  const phat = upvotes / n
  const lower = (phat + z * z / (2 * n) - z * Math.sqrt((phat * (1 - phat) + z * z / (4 * n)) / n)) / (1 + z * z / n)
  // We only want the lower boundary of Wilson's Confidence Interval
  // const upper = (phat + z * z / (2 * n) + z * Math.sqrt((phat * (1 - phat) + z * z / (4 * n)) / n)) / (1 + z * z / n)
  return lower // { lower, upper }
}

console.log(wilson(100, 10))
