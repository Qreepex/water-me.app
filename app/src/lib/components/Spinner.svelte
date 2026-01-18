<div
	class="flex min-h-screen items-center justify-center bg-gradient-to-br from-green-50 via-emerald-50 to-teal-50"
>
	<div class="flex flex-col items-center gap-6">
		<!-- Plant SVG with animation -->
		<div class="relative h-24 w-24">
			<svg
				viewBox="0 0 100 100"
				xmlns="http://www.w3.org/2000/svg"
				class="h-full w-full animate-pulse"
			>
				<!-- Pot -->
				<path
					d="M 25 60 L 30 90 L 70 90 L 75 60 Z"
					fill="#D2691E"
					stroke="#8B4513"
					stroke-width="1.5"
				/>
				<ellipse
					cx="50"
					cy="60"
					rx="25"
					ry="6"
					fill="#CD853F"
					stroke="#8B4513"
					stroke-width="1.5"
				/>

				<!-- Soil -->
				<ellipse cx="50" cy="62" rx="22" ry="4" fill="#A0643F" />

				<!-- Left leaf -->
				<path
					d="M 50 55 Q 35 45 30 25 Q 28 15 35 20 Q 45 32 50 45"
					fill="none"
					stroke="#22C55E"
					stroke-width="3"
					stroke-linecap="round"
					stroke-linejoin="round"
					class="origin-bottom animate-bounce"
					style="animation-delay: 0s; animation-duration: 2s"
				/>

				<!-- Right leaf -->
				<path
					d="M 50 55 Q 65 45 70 25 Q 72 15 65 20 Q 55 32 50 45"
					fill="none"
					stroke="#16A34A"
					stroke-width="3"
					stroke-linecap="round"
					stroke-linejoin="round"
					class="origin-bottom animate-bounce"
					style="animation-delay: 0.3s; animation-duration: 2s"
				/>

				<!-- Center leaf -->
				<path
					d="M 50 55 Q 50 40 50 15"
					fill="none"
					stroke="#22C55E"
					stroke-width="3"
					stroke-linecap="round"
					class="origin-bottom animate-bounce"
					style="animation-delay: 0.6s; animation-duration: 2s"
				/>
			</svg>

			<!-- Spinning ring around the plant -->
			<div
				class="absolute inset-0 animate-spin rounded-full border-4 border-transparent border-t-green-500 border-r-emerald-500"
				style="animation-duration: 3s"
			></div>
		</div>

		<!-- Loading text -->
		<div class="text-center">
			<p class="text-lg font-semibold text-green-700">Growing...</p>
			<p class="mt-1 text-sm text-green-600">Please wait while we load your garden</p>
		</div>

		<!-- Animated dots -->
		<div class="flex gap-1">
			<div
				class="h-2 w-2 animate-bounce rounded-full bg-green-500"
				style="animation-delay: 0s"
			></div>
			<div
				class="h-2 w-2 animate-bounce rounded-full bg-emerald-500"
				style="animation-delay: 0.2s"
			></div>
			<div
				class="h-2 w-2 animate-bounce rounded-full bg-teal-500"
				style="animation-delay: 0.4s"
			></div>
		</div>
	</div>
</div>
