package messages

const (
	HELP_MESSAGE = "راهنمای استفاده از ربات\n\n" +
		"1️⃣ مناقصه عادی همان reverse auction است که در آن نمره‌ی اولیه و حداقل اختلاف اعلام می‌شود؛ سپس افراد پیشنهاد خود را ارسال می‌کنند و وقتی کسی دیگر پیشنهاد جدیدی وارد نکند، پیشنهاد آخر برنده مناقصه می‌شود.\n" +
		"برای شروع کردن چنین مناقصه‌ای به سه ورودی نام مناقصه (شامل حروف و ارقام)، نمره اولیه (float) و حداقل اختلاف پیشنهاد (float) احتیاج دارد. مثال:\n" +
		"/start reverse_auction hw1 3.0 0.02\n" +
		"همچنین برای شرکت در مناقصه باید از دستور /bid استفاده کنید و در ادامه نام مناقصه و نمره پیشنهادی خود را وارد کنید. مثال:\n" +
		"/bid hw1 2.9\n\n" +
		"2️⃣ مناقصه خاص به این صورت است که از یک نمره کم شروع می‌شود و ذره ذره افزایش می‌یابد. اولین کسی که با توجه به نمره فعلی پیشنهاد را قبول کند، برنده مناقصه خواهد بود.\n" +
		"این نوع مناقصه نیز به سه ورودی نام مناقصه (شامل حروف و ارقام)، نمره اولیه (float) و مقدار افزایش در هر گام (float) احتیاج دارد. مثال:\n" +
		"/start special_auction hw2 0 0.02\n\n" +
		"برای شرکت در مناقصه باید از دستور /bid استفاده کنید و در ادامه نام مناقصه را وارد کنید. مثال:\n" +
		"/bid hw2\n"

	START_REVERSE_AUCTION_MESSAGE = " مناقصه جدید %s شروع شد.\n\n" +
		"در این مناقصه از نمره %f شروع می‌کنیم. در هر لحظه، هر کس می‌تواند پیشنهاد جدیدی ثبت کند، پیشنهاداتی که از آخرین پیشنهاد پذیرفته شده حداقل %f کمتر باشد، پذیرفته می‌شوند.\n" +
		"اگر به مدت یک دقیقه هیچ پیشنهادی ثبت نشود، پیشنهاد آخر برنده مناقصه خواهد بود.\n" +
		"برای شرکت در مناقصه باید از دستور /bid استفاده کنید و بعد از آن نام مناقصه و نمره پیشنهادی خود را وارد کنید. مثال:\n" +
		"/bid %s 2.9\n\n" +
		"موفق باشید! 🍀"

	START_SPECIAL_AUCTION_MESSAGE = " مناقصه خاص جدید %s شروع شد.\n\n" +
		"در این مناقصه از نمره %f شروع می‌کنیم و هر ۱۵ ثانیه ارزش سوال %f افزایش می‌یابد.\n" +
		"اولین کسی که پیشنهاد را قبول کند، برنده مناقصه با همان نمره خواهد بود.\n" +
		"حواستون باشه که در این نوع از مناقصه شما هیچ عددی را وارد نمی‌کنید. در واقع کافی است که پیام زیر را ارسال کنید تا پیشنهادتان پذیرفته شود.\n" +
		"/bid %s\n\n" +
		"موفق باشید! 🍀"

	COUNTDOWN_THREE_MESSAGE = "3️⃣ پیشنهاد %f نمره برای این سوال داریم! \n" +
		"بشتابید و نمره کمتری پیشنهاد دهید."

	COUNTDOWN_TWO_MESSAGE = "2️⃣ پیشنهاد %f نمره برای این سوال \n" +
		"کسی نمی‌خواهد پیشنهاد کمتری دهد؟"

	COUNTDOWN_ONE_MESSAGE = "1️⃣ پیشنهاد %f نمره برای این سوال \n" +
		"آخرین فرصت برای پیشهاد کمتر رو از دست نده!"

	NO_ACTIVE_AUCTION_MESSAGE = "هیچ مناقصه فعالی وجود ندارد."

	NO_AUCTION_MESSAGE = "مناقصه با این نام وجود ندارد."

	INVALID_CHAT_ID_MESSAGE = "اینجا جای مناسبی برای پیشنهاد دادن به این مناقصه نیست!"

	INVALID_BID_MESSAGE = "پیشنهاد نامعتبر است. پیشنهاد خود را با فرمت درست ارسال کنید."

	INVALID_BID_AMOUNT_MESSAGE = "همین الانش پیشنهاد %f داریم! پیشنهادت باید کمتر از %f باشه!"

	ACCEPTED_BID_MESSAGE = "*****************************\n" +
		"پیشنهاد %s با نمره %f با موفقیت پذیرفته شد! فعلا شما برنده این مناقصه هستید.\n" +
		"*****************************"

	END_AUCTION_MESSAGE = "مناقصه %s به پایان رسید. 🔨\n" +
		"این سوال زیبا به %s با نمره %f داده شد. به ایشان تبریک می‌گوییم."

	NOT_ADMIN_MESSAGE = "🚫 متاسفانه امکان اینکه شما مناقصه را آغاز کنید وجود ندارد."

	ACTIVE_AUCTION_EXISTS_MESSAGE = "‼️ اول اجازه بده مناقصه قبلی به پایان برسه!"

	INVALID_START_MESSAGE = "❗️ فرمت شروع یک مناقصه را درست وارد کنید. با استفاده از /help می‌توانید فرمت صحیح را پیدا کنید."

	SPECIAL_AUCTION_PRICE_RAISED_MESSAGE = "نمره سوال افزایش یافت و ارزش فعلی سوال %f است.\n"

	SPECIAL_AUCTION_BID_ACCEPTED_MESSAGE = "پیشنهاد %s با موفقیت پذیرفته شد!\n"

	BID_ALREADY_PLACED_MESSAGE = "شرمنده، فقط یک بار می‌توانید پیشنهاد دهید."

	SEALED_BID_ACCEPTED_MESSAGE = "پیشنهاد شما با موفقیت ثبت شد! 📝" +
		"اگر پیشنهاد شما برنده مناقصه شود، از طریق گروه به شما اطلاع داده خواهد شد."

	SEND_BID_PRIVATE_MESSAGE = "اینجا که همه پیشنهادت را می‌بینند! 🤫\n" +
		"برای شرکت در مناقصه، پیشنهاد خود را در چت شخصی با ربات ارسال کنید."

	START_SEALED_BID_AUCTION_MESSAGE = " مناقصه مخفی جدید %s شروع شد.\n\n" +
		"در این مناقصه هر کس تنها می‌تواند یک بار پیشنهاد خود را ثبت کند. پیشنهادات مخفی بوده و تا زمان پایان مناقصه هیچ کس نمی‌تواند بفهمد چه کسی چه مقداری پیشنهاد داده است.\n" +
		"برای شرکت در مناقصه در چت شخصی با ربات پیام دهید. پیام شما باید با دستور /bid و نام مناقصه و مقدار پیشنهادی شما آغاز شود. مثال:\n" +
		"/bid %s 2.9\n\n" +
		"موفق باشید! 🍀"

	SEALED_HALF_TIME_MESSAGE = "نصف زمان مناقصه گذشته است. فرصت شرکت در مناقصه را از دست ندهید!"
)
