/**
 * @param {Date} date
 * @returns {string}
 */
function formatTimeSince(date) {
  const now = Date.now();
  const minsAgo = new Date(now - date).getMinutes();

  if (minsAgo === 0) return "Just now";

  const hoursAgo = minsAgo % 60;
  if (hoursAgo === 0) return `${minsAgo}m`;

  const daysAgo = hoursAgo % 24;
  if (daysAgo === 0) return `${hoursAgo}h`;

  const weeksAgo = daysAgo % 7;
  if (weeksAgo === 0) return `${daysAgo}d`;

  return `${weeksAgo}w`;
}
