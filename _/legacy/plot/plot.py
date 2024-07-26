import matplotlib.pyplot as plt
import numpy as np

plt.rcParams.update({
    "text.usetex": False,
    "font.family": "STIXGeneral",
    "mathtext.fontset": "stix"  
})

keys = [5000, 10000, 20000, 40000, 80000]
dksap = [2362, 4694, 9393, 18797, 37558] 
dksap_view_tag = [308, 614, 1225, 2450, 4885]  # BASE SAP project timings
ecpdksap = [296, 588, 1182, 2355, 4685]  # Double key timings
ecpsksap = [3411, 6813, 13624, 27272, 54435] # Single key timings

# Set positions for the bars
pos = np.arange(len(keys))

# Create the plotting area 
fig, ax = plt.subplots(figsize=(8, 6))

bar_width = 0.2  # Adjust the width of the bars

# Aligning the bars in the center
adjusted_pos = pos - bar_width

# Create bars
rects1 = ax.bar(adjusted_pos, dksap, bar_width, label='DKSAP', align='edge')
rects2 = ax.bar(adjusted_pos + bar_width, dksap_view_tag, bar_width, label='DKSAP with view_tag', align='edge')
rects3 = ax.bar(adjusted_pos + 2 * bar_width, ecpdksap, bar_width, label='ECPDKSAP with view_tag', align='edge')
rects4 = ax.bar(adjusted_pos + 3 * bar_width, ecpsksap, bar_width, label='ECPSKSAP with view_tag', align='edge')

# Add xticks on the middle of the group bars
ax.set_xticks(adjusted_pos + bar_width)
ax.set_xticklabels([str(k) for k in keys])

# Labeling
ax.set_xlabel('Number of Announcements')
ax.set_ylabel('Milliseconds (ms)')
ax.set_title('Performance Comparison')
ax.legend()

def autolabel(rects):
    for rect in rects:
        height = rect.get_height()
        ax.annotate(f'{height}',
                    xy=(rect.get_x() + rect.get_width() / 2, height),
                    xytext=(0, 3),  # 3 points vertical offset
                    textcoords="offset points",
                    ha='center', va='bottom')

autolabel(rects1)
autolabel(rects2)
autolabel(rects3)
autolabel(rects4)

fig.tight_layout()
plt.show()