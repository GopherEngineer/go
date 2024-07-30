import React, { useEffect, useState } from 'react'
import { createRoot } from 'react-dom/client'
import moment from 'moment'
import _ from 'lodash'

import './styles.css'

interface INewsFullDetailed {
	news: IPostDetailed[]
	page: number
	pages: number
	count: number
}

interface IComment {
	id: number
	post: number
	parent: number
	message: string
	created_at: number
}

interface INewsShortDetailed {
	post: IPostDetailed
	comments: IComment[]
}

interface IPostDetailed {
	id: number
	title: string
	content: string
	pub_time: number
	link: string
}

const fetcher = <T,>(url: string, signal: AbortSignal | undefined) => {
	return fetch(url, { signal }).then((res) => res.json()) as Promise<T>
}

interface IPost {
	id: number
	setPostId: (id: number) => void
}

const Post: React.FC<IPost> = ({ id, setPostId }) => {

	const [post, setPost] = useState<IPostDetailed>()
	const [comments, setComments] = useState<IComment[]>()

	const [comment, setComment] = useState<string>('')

	const fetchPost = () => {
		const controller = new AbortController()
		fetcher<INewsShortDetailed>(`http://localhost:8080/news/${id}`, controller.signal).then(({ post, comments }) => {
			setPost(post)
			setComments(comments)
		})
		return () => {
			controller.abort()
		}
	}

	useEffect(() => {
		fetchPost()
	}, [])

	const handleComment = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
		setComment(e.target.value)
	}

	const handleCreateComment = () => {
		fetch(`http://localhost:8080/comments/${id}`, { method: 'post', body: JSON.stringify({ message: comment }) }).then(() => {
			setComment('')
			fetchPost()
		})
	}

	return (
		<div className="flex flex-col items-center h-full mx-auto max-w-[1200px] p-4 gap-6">
			<h1 className="text-6xl max-xl:text-5xl text-center text-[#b5e916] m-2">Новостной агрегатор</h1>

			<button className='bg-white p-2 rounded-lg' onClick={() => setPostId(0)}>Назад к публикациям</button>

			{post?.id && (
				<div
					className="flex flex-col w-full bg-slate-100 hover:bg-white border-4 border-[#328818] hover:border-white transition-all rounded-xl p-4 gap-4 cursor-pointer"
				>
					<div className="flex flex-col flex-1 gap-4">
						<div className="flex justify-between">
							<a className="text-lg text-blue-600 underline" href={post.link}>
								{post.title}
							</a>
							<img className="w-6 h-6 ml-2" src="/img/rss.png" alt={post.title} />
						</div>
						<p className="break-words top line-clamp-6 text-gray-700">{post.content}</p>
					</div>
					<p className="break-words top text-right text-gray-500">{moment(post.pub_time * 1000).format('HH:mm DD-MM-YYYY')}</p>
				</div>
			)}

			<div className='flex flex-col w-full gap-4'>
				<h2 className="text-2xl text-center text-[#b5e916]">Комментарии к публикации</h2>
				{comments?.length > 0 && (
					<div className="flex flex-col gap-2">
						{comments.map((comment) => (
							<div
								key={comment.id}
								className="flex flex-col bg-slate-100 hover:bg-white border-4 border-[#328818] hover:border-white transition-all rounded-xl p-4 gap-4 cursor-pointer"
							>
								<div className="flex flex-col flex-1 gap-4">
									<p className="break-words top line-clamp-6 text-gray-700">{comment.message}</p>
								</div>
								<p className="break-words top text-right text-gray-500">{moment(comment.created_at * 1000).format('HH:mm DD-MM-YYYY')}</p>
							</div>
						))}
					</div>
				)}
				<div className=''>
					<textarea className='p-3 w-full outline-none resize-none rounded-md' placeholder='Ваш комментарий ...' rows={5} value={comment} onChange={handleComment}></textarea>
					<button className='bg-white p-2 rounded-sm' onClick={handleCreateComment}>Сохранить комментарий</button>
				</div>
			</div>
		</div>
	)
}

interface IAggregator {
	setPostId: (id: number) => void
}

const Aggregator: React.FC<IAggregator> = ({ setPostId }) => {
	const [page, setPage] = useState<number>(1)
	const [search, setSearch] = useState<string>('')

	const [pages, setPages] = useState<number>(0)
	const [news, setNews] = useState<IPostDetailed[]>([])

	useEffect(() => {
		const controller = new AbortController()
		fetcher<INewsFullDetailed>(`http://localhost:8080/news?page=${page}&s=${search}`, controller.signal).then(({ news, pages }) => {
			setNews(news)
			setPages(pages)
		})
		return () => {
			controller.abort()
		}
	}, [page, search])

	const handleSearch = _.debounce((e: React.ChangeEvent<HTMLInputElement>) => {
		setSearch(e.target.value)
		setPage(1)
	}, 500)

	return (
		<div className="flex flex-col items-center h-full p-4 gap-6">
			<h1 className="text-6xl max-xl:text-5xl text-center text-[#b5e916] m-2">Новостной агрегатор</h1>
			<input type='text' className='w-[500px] outline-none p-2 rounded-lg' placeholder='Поиск ...' onChange={handleSearch} />
			{pages > 0 && (
				<div className='flex gap-2'>
					{_.range(pages).map((v) => (
						<span key={v} className={`flex justify-center items-center w-8 h-8 text-xl rounded-md cursor-pointer border ${page == v + 1 ? 'bg-white' : 'text-white'}`} onClick={() => setPage(v + 1)}>{v + 1}</span>
					))}
				</div>
			)}
			<div className="grid grid-cols-4 max-sm:grid-cols-1 max-lg:grid-cols-2 max-xl:grid-cols-3 gap-6">
				{news.map((post) => (
					<div
						key={post.id}
						className="flex flex-col bg-slate-100 hover:bg-white border-4 border-[#328818] hover:border-white transition-all rounded-xl p-4 gap-4 cursor-pointer"
						onClick={() => setPostId(post.id)}
					>
						<div className="flex flex-col flex-1 gap-4">
							<div className="flex justify-between">
								<a className="text-lg text-blue-600 underline" href={post.link}>
									{post.title}
								</a>
								<img className="w-6 h-6 ml-2" src="/img/rss.png" alt={post.title} />
							</div>
							<p className="break-words top line-clamp-6 text-gray-700">{post.content}</p>
						</div>
						<p className="break-words top text-right text-gray-500">{moment(post.pub_time * 1000).format('HH:mm DD-MM-YYYY')}</p>
					</div>
				))}
			</div>
		</div>
	)
}

function Main() {
	const [postId, setPostId] = useState<number>(0)
	return postId > 0 ? <Post id={postId} setPostId={setPostId} /> : <Aggregator setPostId={setPostId} />
}

const root = createRoot(document.getElementById('root'))
root.render(<Main />)
